from oss.oss_api import *
endpoint="localhost:9999"
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
bucket = "bucket-name"
object = "object/name"

lifecycle_xml = "<LifecycleConfiguration><Rule><ID>delete obsoleted files</ID><Prefix>obsoleted/</Prefix><Status>Enabled</Status><Expiration><Days>3</Days></Expiration></Rule><Rule><ID>delete temporary files</ID><Prefix>temporary/</Prefix><Status>Enabled</Status><Expiration><Date>2022-10-12T00:00:00.001Z</Date></Expiration></Rule></LifecycleConfiguration>"
res=oss.put_lifecycle(bucket,lifecycle_xml)
