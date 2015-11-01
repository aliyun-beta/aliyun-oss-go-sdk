from oss.oss_api import *
from oss import oss_xml_handler
endpoint="127.0.0.1:9999"
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)

res= oss.get_bucket_location("bucket-name")
acl_xml=oss_xml_handler.GetBucketAclXml(res.read())
print acl_xml.grant
