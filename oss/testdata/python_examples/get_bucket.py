from oss.oss_api import *
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
oss.debug = True
bucket = "bucket-name"

oss.get_bucket(bucket,"private")
