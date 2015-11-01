from oss.oss_api import *
endpoint="localhost:9999"
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
bucket = "bucket-name"
object = "object/name"

res = oss.put_object_from_file(bucket, object, "testdata/test")
print "%s\n%s" % (res.status, res.getheaders())
