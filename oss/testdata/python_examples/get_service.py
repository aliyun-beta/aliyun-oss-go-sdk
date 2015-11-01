from oss.oss_api import *
from oss import oss_xml_handler
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
oss.debug = True

res = oss.get_service()
buckets_xml=oss_xml_handler.GetServiceXml(res.read())
for bucket_info in buckets_xml.bucket_list:
    print "----------------------------------------"
    print "Location:"+bucket_info.location
    print "Name:"+bucket_info.name
    print "CreationDate:"+bucket_info.creation_date
