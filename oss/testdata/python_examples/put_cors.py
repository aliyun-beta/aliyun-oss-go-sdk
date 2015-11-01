from oss.oss_api import *
endpoint="localhost:9999"
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
bucket = "bucket-name"
object = "object/name"

cors_xml = "<CORSConfiguration><CORSRule><AllowedOrigin>*</AllowedOrigin><AllowedMethod>PUT</AllowedMethod><AllowedMethod>GET</AllowedMethod><AllowedHeader>Authorization</AllowedHeader></CORSRule><CORSRule><AllowedOrigin>http://www.a.com</AllowedOrigin><AllowedOrigin>http://www.b.com</AllowedOrigin><AllowedMethod>GET</AllowedMethod><AllowedHeader>Authorization</AllowedHeader><ExposeHeader>x-oss-test</ExposeHeader><ExposeHeader>x-oss-test1</ExposeHeader><MaxAgeSeconds>100</MaxAgeSeconds></CORSRule></CORSConfiguration>"
res=oss.put_cors(bucket,cors_xml)
