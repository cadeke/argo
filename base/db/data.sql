-- Insert servers with IP addresses
INSERT INTO servers (name, ip, contact, company) VALUES
('Google DNS', '8.8.8.8', 'admin@google.com', 'Google LLC'),
('Cloudflare DNS', '1.1.1.1', 'support@cloudflare.com', 'Cloudflare Inc.'),
('OpenDNS', '208.67.222.222', 'help@opendns.com', 'Cisco Systems'),
('Quad9 DNS', '9.9.9.9', 'security@quad9.net', 'Quad9 Foundation'),
('Custom DNS', '192.168.1.1', 'contact@customdns.net', 'Custom Networks');


INSERT INTO records (domain, ip, server) VALUES
-- Google DNS (Server ID 1)
('example.com', '93.184.216.34', 1),
('google.com', '8.8.8.8', 1),
('youtube.com', '142.250.190.46', 1),
('gmail.com', '142.251.33.229', 1),
('maps.google.com', '142.250.185.78', 1),
('news.google.com', '142.250.188.206', 1),
('drive.google.com', '142.250.187.78', 1),
('play.google.com', '142.250.180.78', 1),
('meet.google.com', '172.217.10.238', 1),
('photos.google.com', '172.217.9.206', 1),

-- Cloudflare DNS (Server ID 2)
('cloudflare.com', '104.16.132.229', 2),
('1.1.1.1', '1.1.1.1', 2),
('example.net', '93.184.216.34', 2),
('wikipedia.org', '208.80.154.224', 2),
('discord.com', '162.159.135.232', 2),
('medium.com', '162.159.152.4', 2),
('notion.so', '104.248.78.23', 2),
('github.com', '140.82.114.3', 2),
('twitter.com', '104.244.42.129', 2),
('facebook.com', '157.240.22.35', 2),

-- OpenDNS (Server ID 3)
('opendns.com', '208.67.222.222', 3),
('linkedin.com', '108.174.10.10', 3),
('stackoverflow.com', '151.101.129.69', 3),
('yahoo.com', '98.138.219.231', 3),
('amazon.com', '176.32.103.205', 3),
('bing.com', '204.79.197.200', 3),
('duckduckgo.com', '52.250.42.157', 3),
('netflix.com', '52.89.124.206', 3),
('hulu.com', '23.246.0.5', 3),
('apple.com', '17.253.144.10', 3),

-- Quad9 DNS (Server ID 4)
('quad9.net', '9.9.9.9', 4),
('microsoft.com', '104.215.148.63', 4),
('adobe.com', '192.147.130.1', 4),
('espn.com', '199.181.132.250', 4),
('cnn.com', '151.101.193.67', 4),
('bbc.com', '151.101.128.81', 4),
('reddit.com', '151.101.65.140', 4),
('pinterest.com', '151.101.64.84', 4),
('tiktok.com', '23.40.242.42', 4),
('spotify.com', '35.186.224.25', 4),

-- Custom DNS (Server ID 5)
('customdns.net', '192.168.1.1', 5),
('localnetwork.com', '192.168.0.1', 5),
('mywebsite.com', '203.0.113.1', 5),
('testdomain.com', '198.51.100.1', 5),
('exampledomain.net', '192.0.2.1', 5),
('fictionalsite.com', '203.0.113.45', 5),
('randompage.net', '192.88.99.1', 5),
('privatehost.net', '192.18.43.1', 5),
('experimental.com', '203.0.113.99', 5),
('unusedsite.com', '198.51.100.99', 5);
