from oss.oss_api import *
from oss import oss_xml_handler
endpoint="127.0.0.1:9999"
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
oss.debug=True
bucket="bucket-name"
object="object/name"
filename="testdata/test"


upload_id="0004B9895DBBB6EC98E36"

res=oss.get_all_parts(bucket,object,upload_id)
