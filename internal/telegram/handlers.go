package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/Anton-Kraev/gopay"
)

func (t *Telegram) handleUpdate(ctx context.Context, update telego.Update) {
	if update.Message == nil {
		return
	}

	var (
		chatID = update.Message.Chat.ID
		err    error
	)

	if !t.checkAuth(chatID) {
		err = t.handleUnauthenticated(ctx, update)
	} else {
		err = t.handleMessage(ctx, update)
	}

	if err != nil {
		t.log.With(
			slog.Int64("chat_id", chatID),
		).Error(
			fmt.Errorf("telegram.handleUpdate: %w", err).Error(),
		)
	}
}

func (t *Telegram) checkAuth(chatID int64) bool {
	return slices.Contains(t.whitelist, chatID)
}

func (t *Telegram) handleUnauthenticated(ctx context.Context, update telego.Update) error {
	return t.sendMessage(
		ctx,
		update,
		"telegram.handleUnauthenticated",
		"у вас нет прав доступа к функциям бота, обратитесь к администратору",
	)
}

func (t *Telegram) handleMessage(ctx context.Context, update telego.Update) error {
	var err error

	text := strings.Split(update.Message.Text, " ")
	if len(text) == 0 {
		return t.sendMessage(
			ctx,
			update,
			"telegram.handleMessage",
			"неверный формат сообщения: ожидается текст",
		)
	}

	switch text[0] {
	case cmdStart:
		err = t.handleCmdStart(ctx, update)
	case cmdNewPayment:
		err = t.handleCmdNewPayment(ctx, update)
	case cmdAllPayment:
		err = t.handleCmdAllPayment(ctx, update)
	case cmdGetPayment:
		err = t.handleCmdGetPayment(ctx, update)
	default:
		err = t.handleState(ctx, update)
	}

	if err != nil {
		return fmt.Errorf("telegram.handleMessage: %w", err)
	}

	return nil
}

func (t *Telegram) handleCmdStart(ctx context.Context, update telego.Update) error {
	delete(t.fsm, update.Message.Chat.ID)

	return t.sendMessage(
		ctx,
		update,
		"telegram.handleCmdStart",
		`вы успешно авторизованы, список доступных команд:
				1) /new_payment --- создание нового платежа
				2) /all_payment --- получение статусов всех платежей
				3) /get_payment <id> --- получение статуса платежа по его id
			`,
	)
}

func (t *Telegram) handleCmdNewPayment(ctx context.Context, update telego.Update) error {
	chatID := update.Message.Chat.ID
	t.newPaymentService[chatID] = t.adminClient.NewNewPaymentService()
	t.fsm[chatID] = stateNewPaymentAmount

	return t.sendMessage(
		ctx,
		update,
		"telegram.handleCmdNewPayment",
		"для создания платежа введите сумму и валюту в формате \"1000 RUB\"",
	)
}

func (t *Telegram) handleCmdAllPayment(ctx context.Context, update telego.Update) error {
	delete(t.fsm, update.Message.Chat.ID)

	statuses, err := t.adminClient.NewAllPaymentService().Do()
	if err != nil {
		return errors.Join(
			fmt.Errorf("telegram.handleCmdAllPayment: %w", err),
			t.sendMessage(
				ctx,
				update,
				"telegram.handleCmdAllPayment",
				"не удалось получить статусы платежей",
			),
		)
	}

	msg := strings.Builder{}
	msg.WriteString("список статусов платежей в формате \"id: status\"")

	for id, status := range statuses {
		msg.WriteString(fmt.Sprintf("\n%s: %s", id, status))
	}

	return t.sendMessage(ctx, update, "telegram.handleCmdAllPayment", msg.String())
}

func (t *Telegram) handleCmdGetPayment(ctx context.Context, update telego.Update) error {
	delete(t.fsm, update.Message.Chat.ID)

	text := strings.Split(update.Message.Text, " ")
	if len(text) != 2 {
		return t.sendMessage(
			ctx,
			update,
			"telegram.handleCmdGetPayment",
			"неверный формат команды, ожидается \"/get_payment <id>\"",
		)
	}

	id := text[1]
	status, err := t.adminClient.NewGetPaymentService().ID(gopay.ID(id)).Do()
	if err != nil {
		return errors.Join(
			fmt.Errorf("telegram.handleCmdGetPayment: %w", err),
			t.sendMessage(
				ctx,
				update,
				"telegram.handleCmdGetPayment",
				"не удалось получить статус платежа "+id,
			),
		)
	}

	return t.sendMessage(
		ctx,
		update,
		"telegram.handleCmdGetPayment",
		fmt.Sprintf("статус платежа %s: %s", id, status),
	)
}

func (t *Telegram) handleState(ctx context.Context, update telego.Update) error {
	var err error

	switch t.fsm[update.Message.Chat.ID] {
	case stateNewPaymentAmount:
		err = t.handleStateNewPaymentAmount(ctx, update)
	case stateNewPaymentDescription:
		err = t.handleStateNewPaymentDescription(ctx, update)
	case stateNewPaymentLink:
		err = t.handleStateNewPaymentLink(ctx, update)
	case stateNewPaymentConfirmation:
		err = t.handleStateNewPaymentConfirmation(ctx, update)
	default:
		err = t.handleUnknownMessage(ctx, update)
	}

	if err != nil {
		return fmt.Errorf("telegram.handleState: %w", err)
	}

	return nil
}

func (t *Telegram) handleStateNewPaymentAmount(ctx context.Context, update telego.Update) error {
	text := strings.Split(update.Message.Text, " ")
	if len(text) != 2 {
		return t.sendMessage(
			ctx,
			update,
			"telegram.handleStateNewPaymentAmount",
			"неверный формат суммы и валюты платежа, пример \"1000 RUB\"",
		)
	}

	amount, err := strconv.ParseUint(text[0], 10, 32)
	if err != nil || amount <= 0 {
		return t.sendMessage(
			ctx,
			update,
			"telegram.handleStateNewPaymentAmount",
			"некорректная сумма платежа, ожидается положительное целое число",
		)
	}

	currency := text[1]
	if currency != "RUB" {
		return t.sendMessage(
			ctx,
			update,
			"telegram.handleStateNewPaymentAmount",
			"некорректная валюта платежа, пример \"RUB\"",
		)
	}

	// TODO: check currency more correctly, currency stdlib locale package

	chatID := update.Message.Chat.ID
	t.newPaymentService[chatID].Currency(currency).Amount(uint(amount))
	t.fsm[update.Message.Chat.ID] = stateNewPaymentDescription

	return t.sendMessage(
		ctx,
		update,
		"telegram.handleStateNewPaymentAmount",
		"сумма и валюта платежа успешно добавлены\nвведите описание платежа:",
	)
}

func (t *Telegram) handleStateNewPaymentDescription(ctx context.Context, update telego.Update) error {
	description := update.Message.Text
	// TODO: maybe some checks here

	chatID := update.Message.Chat.ID
	t.newPaymentService[chatID].Description(description)
	t.fsm[chatID] = stateNewPaymentLink

	return t.sendMessage(
		ctx,
		update,
		"telegram.handleStateNewPaymentDescription",
		"описание платежа успешно добавлено\nвведите ссылку на ресурс:",
	)
}

func (t *Telegram) handleStateNewPaymentLink(ctx context.Context, update telego.Update) error {
	link := gopay.Link(update.Message.Text)
	if !link.Validate() {
		return t.sendMessage(
			ctx,
			update,
			"telegram.handleStateNewPaymentLink",
			"некорректная ссылка",
		)
	}

	chatID := update.Message.Chat.ID
	t.newPaymentService[chatID].ResourceLink(link)
	t.fsm[chatID] = stateNewPaymentConfirmation

	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("да"),
			tu.KeyboardButton("нет"),
		),
	)

	msg := tu.Message(
		tu.ID(chatID),
		"ссылка на ресурс успешно добавлена, подтвердить создание платежа со следующими параметрами?:\n"+
			t.newPaymentService[chatID].String(),
	).WithReplyMarkup(keyboard)

	_, err := t.bot.SendMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("telegram.handleStateNewPaymentDescription: %w", err)
	}

	return nil
}

func (t *Telegram) handleStateNewPaymentConfirmation(ctx context.Context, update telego.Update) error {
	text := strings.Split(update.Message.Text, " ")
	if len(text) == 1 && (text[0] != "да" && text[0] != "нет") {
		return t.sendMessage(
			ctx,
			update,
			"telegram.handleStateNewPaymentConfirmation",
			"некорректный ответ, введите \"да\"/\"нет\" или нажмите на одну из соответствующих кнопок",
		)
	}

	chatID := update.Message.Chat.ID
	defer delete(t.newPaymentService, chatID)
	defer delete(t.fsm, chatID)

	if text[0] == "да" {
		link, err := t.newPaymentService[chatID].Do()
		if err != nil {
			return errors.Join(
				fmt.Errorf("telegram.handleStateNewPaymentConfirmation: %w", err),
				t.sendMessage(
					ctx,
					update,
					"telegram.handleStateNewPaymentConfirmation",
					"не удалось создать платеж",
				),
			)
		}

		return t.sendMessage(
			ctx,
			update,
			"telegram.handleStateNewPaymentConfirmation",
			"платеж успешно создан, платежная ссылка:\n"+string(link),
		)
	}

	return t.sendMessage(
		ctx,
		update,
		"telegram.handleStateNewPaymentConfirmation",
		"создание нового платежа отменено",
	)
}

func (t *Telegram) handleUnknownMessage(ctx context.Context, update telego.Update) error {
	return t.sendMessage(
		ctx,
		update,
		"telegram.handleUnknownMessage",
		"неизвестная команда, введите /start для получения списка команд",
	)
}
