package config

const OTP_TEMPLATE = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>OTP Verification</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f8; font-family:Arial, sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f8; padding:20px 0;">
    <tr>
        <td align="center">
            <!-- Main Container -->
            <table width="400" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; padding:30px; text-align:center; box-shadow:0 2px 8px rgba(0,0,0,0.05);">

                <!-- Logo / App Name -->
                <tr>
                    <td style="font-size:22px; font-weight:bold; color:#333;">
                    Mantritva
                    </td>
                </tr>

                <!-- Spacer -->
                <tr><td height="20"></td></tr>

                <!-- Heading -->
                <tr>
                    <td style="font-size:18px; color:#333;">
                    Your One-Time Password (OTP)
                    </td>
                </tr>

                <!-- Spacer -->
                <tr><td height="10"></td></tr>

                <!-- Description -->
                <tr>
                    <td style="font-size:14px; color:#666;">
                    Use the OTP below to complete account verification. This code is valid for the next <b>5 minutes</b>.
                    </td>
                </tr>

                <!-- Spacer -->
                <tr><td height="25"></td></tr>

                <!-- OTP Box -->
                <tr>
                    <td>
                    <div style="display:inline-block; padding:15px 25px; font-size:28px; letter-spacing:6px; font-weight:bold; color:#2d89ef; background:#f1f7ff; border-radius:6px;">
                        {{otp}}
                    </div>
                    </td>
                </tr>

                <!-- Spacer -->
                <tr><td height="25"></td></tr>

                <!-- Warning -->
                <tr>
                    <td style="font-size:13px; color:#999;">
                    Do not share this OTP with anyone. If you did not request this, you can safely ignore this email.
                    </td>
                </tr>

                <!-- Spacer -->
                <tr><td height="30"></td></tr>

                <!-- Footer -->
                <tr>
                    <td style="font-size:12px; color:#bbb;">
                    © 2026 Mantritva. All rights reserved.
                    </td>
                </tr>

            </table>
        </td>
    </tr>
  </table>
</body>
</html>`
