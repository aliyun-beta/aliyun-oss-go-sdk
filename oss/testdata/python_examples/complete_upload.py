from oss.oss_api import *
from oss import oss_util
from oss import oss_xml_handler
endpoint="127.0.0.1:9999"
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
oss.debug=True
bucket="bucket-name"
object="object/name"
filename="testdata/test"


part_msg_list = oss_util.split_large_file(filename, object, 3)

part_number = 1
partsize = 16
offset = 23
upload_id="0004B9895DBBB6EC98E36"
part_msg_xml = oss_util.create_part_xml(part_msg_list)
res = oss.complete_upload(bucket,object, upload_id, "<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>C1B61751512FFC8B0E86675D114497A6</ETag></Part></CompleteMultipartUpload>")
