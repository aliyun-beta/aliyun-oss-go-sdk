from oss.oss_api import *
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
bucket = "bucket-name"

oss.put_bucket(bucket,"private")
print "%s\n%s" % (res.status, res.read())
