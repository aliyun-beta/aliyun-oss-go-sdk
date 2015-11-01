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


part_number = 1
partsize = 16
offset = 23
upload_id="0004B9895DBBB6EC98E36"
res = oss.upload_part_from_file_given_pos(bucket, object, filename, offset, partsize, upload_id, part_number)
