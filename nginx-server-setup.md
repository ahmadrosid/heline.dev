# Nginx Server Setup Tutorial for Heline API

This guide explains how to set up Nginx as a reverse proxy with SSL for the Heline API application.

## Prerequisites

- A server running Ubuntu (or similar Linux distribution)
- Domain name pointed to your server (in this case, api.heline.dev)
- Heline application running on port 8000

## Step 1: Install Nginx

Update the package list and install Nginx:

```bash
apt update
apt install -y nginx
```

## Step 2: Create Nginx Configuration

Create a new configuration file for your domain:

```bash
nano /etc/nginx/sites-available/api.heline.dev
```

Add the following configuration:

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name api.heline.dev;

    location / {
        proxy_pass http://localhost:8000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Step 3: Enable the Site

Remove the default site and enable your configuration:

```bash
rm /etc/nginx/sites-enabled/default
ln -s /etc/nginx/sites-available/api.heline.dev /etc/nginx/sites-enabled/
```

Test the configuration:

```bash
nginx -t
```

If the test is successful, restart Nginx:

```bash
systemctl restart nginx
```

## Step 4: Install and Configure SSL

1. Install Certbot and its Nginx plugin:

```bash
apt install -y certbot python3-certbot-nginx
```

2. Obtain and install SSL certificate:

```bash
certbot --nginx -d api.heline.dev --non-interactive --agree-tos --redirect
```

This command will:
- Obtain an SSL certificate
- Modify your Nginx configuration to use HTTPS
- Set up automatic redirects from HTTP to HTTPS
- Configure automatic renewal

## Verification

After setup, you can verify:

1. Nginx status:
```bash
systemctl status nginx
```

2. Certbot timer status (for auto-renewal):
```bash
systemctl status certbot.timer
```

## Final Configuration

Your site should now be accessible via:
- https://api.heline.dev (secure access)
- http://api.heline.dev (automatically redirects to HTTPS)

## Maintenance

- SSL certificates will automatically renew every 90 days
- The Nginx configuration file is located at `/etc/nginx/sites-available/api.heline.dev`
- Logs can be found in `/var/log/nginx/`

## Additional Notes

- The configuration includes WebSocket support
- Proper security headers are set
- X-Forwarded-* headers are properly configured for proxy setup

## Troubleshooting

If you encounter issues:

1. Check Nginx logs:
```bash
tail -f /var/log/nginx/error.log
```

2. Verify Nginx configuration:
```bash
nginx -t
```

3. Check SSL certificate status:
```bash
certbot certificates
```

4. Ensure firewall allows ports 80 and 443:
```bash
ufw status
```

## SSL Certificate Auto-Renewal

The SSL certificate auto-renewal is handled by the Certbot systemd timer, which is already configured during installation. This is more reliable than a traditional cron job.

### Verify Auto-Renewal Setup

Check the timer status:
```bash
systemctl status certbot.timer
```

The timer is configured to:
- Run twice daily
- Attempt renewal only when certificates are near expiration (within 30 days)
- Automatically reload Nginx after successful renewal

### Manual Renewal Test

To test the renewal process without actually renewing the certificate:
```bash
certbot renew --dry-run
```

### Renewal Logs

Check the renewal logs:
```bash
journalctl -u certbot.service
```

### Important Paths

- Certificates: `/etc/letsencrypt/live/api.heline.dev/`
- Renewal configuration: `/etc/letsencrypt/renewal/api.heline.dev.conf`
- Logs: `/var/log/letsencrypt/`

### Backup Considerations

Consider backing up the following directories:
```bash
/etc/letsencrypt/
/var/lib/letsencrypt/
```

These contain your certificates and renewal configuration.

### Monitoring

To ensure the renewal system is working:
1. Check the systemd timer status regularly
2. Monitor certificate expiration dates:
   ```bash
   certbot certificates
   ```
3. Set up monitoring for the HTTPS certificate expiry date
4. Configure alerts if renewal fails

The current setup will attempt renewal twice daily when the certificate is within 30 days of expiration. The renewal process is fully automated and requires no manual intervention under normal circumstances.
