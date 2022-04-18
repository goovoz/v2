import hashlib
import urllib.parse

s = urllib.parse.urlparse('https://elder.kugou.com/v1/incentive/user_info?appid=3176&token=h5747D8FB94D1964A20CB059C9E418FC4349693B7E024CC2123868F55D5E919106DEB0B3D176567AF6805B8056912763BA0EE03656FA6F00C85D917795D1719724C51BF82056600FDABA4EB6B143800AECBEA768BB78E8FF67B0B92D18F3F2044713F3AB0CA309D17574571373F132202AF5E54B6E095DF149C2E7DAE06D4B1608&userid=1413252851&clienttime=1650459646&dfid=4WQcNn1Oj7fr0XJHoQ1eucEt&mid=0af86bb3eafb054d7aab9f984a7760aabe36db29&signature=225d291cf2b8be66e6a1c6f5331953f4&clientver=12170')

t = []

for item in s.query.split('&'):
    if 'signature' in item:
        continue
    t.append(item)

import hmac


t = sorted(t, reverse=True)


s2 = '&'.join(t)
key = '4WQcNn1Oj7fr0XJHoQ1eucEt'
h_mac = hmac.new(key.encode(), s2.encode(), digestmod = 'sha1')
print(h_mac.hexdigest()) #fd34d13d4e31d362f19f1fa9e783fcf0
