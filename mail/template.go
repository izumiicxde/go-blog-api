package mail

func MailTemplate(verificationCode string, username string) string {
	return `
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>NAX Blogs</title>
			<style>
				body {
					font-family: Helvetica, Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 20px;
					color: #333;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #4CAF50;
					text-align: center;
				}
				p {
					font-size: 16px;
					line-height: 1.5;
					text-align: center;
				}
				h2 {
					color: #333;
					text-align: center;
					font-size: 24px;
					margin: 20px 0;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Nax blogs | katana</h1>
				<p style="font-size: 32px;">Hello ` + username + `,</p>
				<p>You have requested a verification code for your Nax blogs account.</p>
				<h2>Your verification code is: ` + verificationCode + `</h2>
				<p>If you did not request this, please ignore this email.</p>
				<p>Thank you for using Nax blogs.</p>
			</div>
		</body>	
	</html>
`
}
