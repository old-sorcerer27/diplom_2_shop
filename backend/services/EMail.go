package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-mail/mail/v2"
)

type OrderItem struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type OrderData struct {
	ID            uint        `json:"id"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	Phone         string      `json:"phone"`
	Address       string      `json:"address"`
	Total         float64     `json:"total"`
	Status        string      `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	Items         []OrderItem `json:"items"`
}

type EmailService struct {
	dialer *mail.Dialer
	from   string
	mu     sync.Mutex
	queue  chan *EmailMessage
}

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

func NewEmailService() (*EmailService, error) {
	if AppConfig.SMTPUser == "" || AppConfig.SMTPPassword == "" {
		log.Println("⚠️ Email credentials not configured, email notifications disabled")
		return &EmailService{}, nil
	}

	dialer := mail.NewDialer(
		AppConfig.SMTPHost,
		AppConfig.SMTPPort,
		AppConfig.SMTPUser,
		AppConfig.SMTPPassword,
	)

	service := &EmailService{
		dialer: dialer,
		from:   AppConfig.SMTPFrom,
		queue:  make(chan *EmailMessage, 100),
	}

	// Запускаем обработчик очереди
	go service.processQueue()

	// Проверяем подключение
	if err := service.testConnection(); err != nil {
		log.Printf("⚠️ Email service test failed: %v", err)
	} else {
		log.Println("✅ Email service connected successfully")
	}

	return service, nil
}

func (s *EmailService) testConnection() error {
	if s.dialer == nil {
		return fmt.Errorf("email service not configured")
	}

	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", AppConfig.SMTPUser)
	m.SetHeader("Subject", "Test connection")
	m.SetBody("text/plain", "Test email from store")

	return s.dialer.DialAndSend(m)
}

func (s *EmailService) SendEmail(to []string, subject, body string, isHTML bool) {
	if s.dialer == nil {
		log.Println("⚠️ Email service not configured, skipping send")
		return
	}

	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
	}

	select {
	case s.queue <- msg:
		log.Printf("📧 Email queued: %s -> %v", subject, to)
	default:
		log.Printf("⚠️ Email queue full, sending directly: %s", subject)
		go s.sendImmediate(msg)
	}
}

func (s *EmailService) processQueue() {
	for msg := range s.queue {
		s.sendImmediate(msg)
		time.Sleep(100 * time.Millisecond) // Rate limiting
	}
}

func (s *EmailService) sendImmediate(msg *EmailMessage) {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)

	if msg.IsHTML {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	// Retry logic
	for i := 0; i < 3; i++ {
		if err := s.dialer.DialAndSend(m); err != nil {
			log.Printf("❌ Failed to send email (attempt %d/3): %v", i+1, err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("✅ Email sent successfully: %s", msg.Subject)
		return
	}

	log.Printf("❌ Failed to send email after 3 attempts: %s", msg.Subject)
}

func formatPrice(price float64) string {
	return fmt.Sprintf("%.2f", price)
}

func (s *EmailService) GenerateOwnerEmail(order *OrderData) string {
	var itemsHTML string
	for _, item := range order.Items {
		itemsHTML += fmt.Sprintf(`
            <tr>
                <td style="padding: 12px; border-bottom: 1px solid #eee;">%s</td>
                <td style="padding: 12px; border-bottom: 1px solid #eee; text-align: center;">%d</td>
                <td style="padding: 12px; border-bottom: 1px solid #eee; text-align: right;">%s ₽</td>
                <td style="padding: 12px; border-bottom: 1px solid #eee; text-align: right;">%s ₽</td>
            </tr>`,
			item.ProductName,
			item.Quantity,
			formatPrice(item.Price),
		)
	}

	return fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <style>
                * { margin: 0; padding: 0; box-sizing: border-box; }
                body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #f5f7fb; line-height: 1.6; }
                .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 15px; overflow: hidden; box-shadow: 0 10px 30px rgba(0,0,0,0.1); }
                .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; }
                .header h1 { margin: 0; font-size: 28px; }
                .header p { margin: 10px 0 0; opacity: 0.9; }
                .content { padding: 30px; }
                .order-info { background: #f8f9fa; border-radius: 10px; padding: 15px; margin-bottom: 20px; }
                .order-info h3 { margin-bottom: 10px; color: #667eea; }
                .order-details { display: flex; justify-content: space-between; flex-wrap: wrap; gap: 15px; }
                .detail-item { flex: 1; min-width: 150px; }
                .detail-label { font-size: 12px; color: #999; text-transform: uppercase; margin-bottom: 5px; }
                .detail-value { font-size: 16px; font-weight: 500; color: #333; }
                .badge { display: inline-block; padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: 600; }
                .badge-pending { background: #ffc107; color: #333; }
                table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
                th { background: #f8f9fa; padding: 12px; text-align: left; font-weight: 600; color: #555; border-bottom: 2px solid #e0e0e0; }
                td { padding: 12px; border-bottom: 1px solid #eee; }
                .total-row { background: #f8f9fa; font-weight: bold; }
                .total-row td { padding: 15px; }
                .footer { background: #f8f9fa; padding: 20px; text-align: center; color: #999; font-size: 12px; }
                .btn { display: inline-block; padding: 10px 20px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; text-decoration: none; border-radius: 8px; margin-top: 20px; }
                @media (max-width: 600px) { .order-details { flex-direction: column; } }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>🛒 Новый заказ #%d</h1>
                    <p>Поступил новый заказ в магазине</p>
                </div>
                
                <div class="content">
                    <div class="order-info">
                        <h3>📋 Информация о заказе</h3>
                        <div class="order-details">
                            <div class="detail-item">
                                <div class="detail-label">Номер заказа</div>
                                <div class="detail-value">#%d</div>
                            </div>
                            <div class="detail-item">
                                <div class="detail-label">Статус</div>
                                <div class="detail-value"><span class="badge badge-pending">%s</span></div>
                            </div>
                            <div class="detail-item">
                                <div class="detail-label">Дата</div>
                                <div class="detail-value">%s</div>
                            </div>
                        </div>
                    </div>
                    
                    <div class="order-info">
                        <h3>👤 Информация о покупателе</h3>
                        <div class="order-details">
                            <div class="detail-item">
                                <div class="detail-label">Имя</div>
                                <div class="detail-value">%s</div>
                            </div>
                            <div class="detail-item">
                                <div class="detail-label">Email</div>
                                <div class="detail-value">%s</div>
                            </div>
                            <div class="detail-item">
                                <div class="detail-label">Телефон</div>
                                <div class="detail-value">%s</div>
                            </div>
                            <div class="detail-item">
                                <div class="detail-label">Адрес доставки</div>
                                <div class="detail-value">%s</div>
                            </div>
                        </div>
                    </div>
                    
                    <h3>🛍️ Состав заказа</h3>
                    <table>
                        <thead>
                            <tr>
                                <th>Товар</th>
                                <th style="text-align: center">Кол-во</th>
                                <th style="text-align: right">Цена</th>
                                <th style="text-align: right">Сумма</th>
                            </tr>
                        </thead>
                        <tbody>
                            %s
                            <tr class="total-row">
                                <td colspan="3" style="text-align: right; font-weight: bold;">Итого:</td>
                                <td style="text-align: right; font-weight: bold; font-size: 18px; color: #28a745;">%s ₽</td>
                            </tr>
                        </tbody>
                    </table>
                    
                    <div style="text-align: center;">
                        <a href="%%s/admin/orders/%d" class="btn">📦 Перейти к заказу</a>
                    </div>
                </div>
                
                <div class="footer">
                    <p>Это письмо отправлено автоматически с сайта интернет-магазина.</p>
                    <p>© 2024 Мой Магазин</p>
                </div>
            </div>
        </body>
        </html>`,
		order.ID,
		order.ID,
		order.Status,
		order.CreatedAt.Format("02.01.2006 15:04:05"),
		order.CustomerName,
		order.CustomerEmail,
		order.Phone,
		order.Address,
		itemsHTML,
		formatPrice(order.Total),
		AppConfig.Port,
		order.ID,
	)
}

// Генерация HTML письма для клиента (подтверждение)
func (s *EmailService) GenerateCustomerEmail(order *OrderData) string {
	var itemsHTML string
	for _, item := range order.Items {
		itemsHTML += fmt.Sprintf(`
            <tr>
                <td style="padding: 12px; border-bottom: 1px solid #eee;">%s</td>
                <td style="padding: 12px; border-bottom: 1px solid #eee; text-align: center;">%d</td>
                <td style="padding: 12px; border-bottom: 1px solid #eee; text-align: right;">%s ₽</td>
            </tr>`,
			item.ProductName,
			item.Quantity,
		)
	}

	return fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <style>
                * { margin: 0; padding: 0; box-sizing: border-box; }
                body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #f5f7fb; line-height: 1.6; }
                .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 15px; overflow: hidden; box-shadow: 0 10px 30px rgba(0,0,0,0.1); }
                .header { background: linear-gradient(135deg, #28a745 0%%, #20c997 100%%); color: white; padding: 30px; text-align: center; }
                .header h1 { margin: 0; font-size: 28px; }
                .content { padding: 30px; }
                .info-box { background: #d4edda; border-left: 4px solid #28a745; padding: 15px; margin: 20px 0; border-radius: 8px; }
                table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
                th { background: #f8f9fa; padding: 12px; text-align: left; border-bottom: 2px solid #e0e0e0; }
                td { padding: 12px; border-bottom: 1px solid #eee; }
                .total { text-align: right; font-size: 18px; font-weight: bold; color: #28a745; margin-top: 20px; }
                .footer { background: #f8f9fa; padding: 20px; text-align: center; color: #999; font-size: 12px; }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>✅ Заказ #%d подтвержден</h1>
                    <p>Спасибо за покупку!</p>
                </div>
                
                <div class="content">
                    <p>Здравствуйте, <strong>%s</strong>!</p>
                    <p>Ваш заказ успешно оформлен. Мы свяжемся с вами в ближайшее время.</p>
                    
                    <div class="info-box">
                        <strong>📋 Детали заказа #%d</strong><br>
                        Дата: %s<br>
                        Статус: <span style="color: #ffc107;">Ожидает обработки</span>
                    </div>
                    
                    <h3>Состав заказа:</h3>
                    <table>
                        <thead>
                            <tr><th>Товар</th><th>Кол-во</th><th>Сумма</th></tr>
                        </thead>
                        <tbody>
                            %s
                        </tbody>
                    </table>
                    
                    <div class="total">
                        Итого: %s ₽
                    </div>
                    
                    <p style="margin-top: 20px;">
                        Следить за статусом заказа вы можете в <a href="%%s/profile">личном кабинете</a>.
                    </p>
                </div>
                
                <div class="footer">
                    <p>© 2024 Мой Магазин</p>
                </div>
            </div>
        </body>
        </html>`,
		order.ID,
		order.CustomerName,
		order.ID,
		order.CreatedAt.Format("02.01.2006 15:04:05"),
		itemsHTML,
		formatPrice(order.Total),
		AppConfig.Port,
	)
}

func (s *EmailService) NotifyOwner(order *OrderData) {
	subject := fmt.Sprintf("🛒 Новый заказ #%d на сумму %s ₽", order.ID, formatPrice(order.Total))
	body := s.GenerateOwnerEmail(order)

	s.SendEmail(
		[]string{AppConfig.OwnerEmail},
		subject,
		body,
		true,
	)
}

func (s *EmailService) NotifyCustomer(order *OrderData) {
	subject := fmt.Sprintf("✅ Заказ #%d подтвержден", order.ID)
	body := s.GenerateCustomerEmail(order)

	s.SendEmail(
		[]string{order.CustomerEmail},
		subject,
		body,
		true,
	)
}

func (s *EmailService) NotifyStatusChange(order *OrderData, oldStatus, newStatus string) {
	statusMessages := map[string]string{
		"paid":      "оплачен и передан в обработку",
		"shipped":   "отправлен. Отслеживайте его в личном кабинете",
		"delivered": "доставлен. Спасибо, что выбрали нас!",
		"cancelled": "отменен. Если это ошибка, свяжитесь с поддержкой",
	}

	message, ok := statusMessages[newStatus]
	if !ok {
		return
	}

	body := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head><meta charset="UTF-8"></head>
        <body style="font-family: Arial, sans-serif;">
            <h2>Статус заказа #%d изменен</h2>
            <p>Здравствуйте, <strong>%s</strong>!</p>
            <p>Ваш заказ <strong>%s</strong>.</p>
            <p>Текущий статус: <strong style="color: #28a745;">%s</strong></p>
            <a href="http://localhost:8080/profile/orders/%d">Перейти к заказу</a>
        </body>
        </html>`,
		order.ID,
		order.CustomerName,
		message,
		newStatus,
		order.ID,
	)

	s.SendEmail(
		[]string{order.CustomerEmail},
		fmt.Sprintf("🔄 Статус заказа #%d: %s", order.ID, newStatus),
		body,
		true,
	)
}
