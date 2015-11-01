from oss.oss_api import *
from oss import oss_xml_handler
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
oss.debug=True
bucket="bucket-name"
object="object/name"
localfile="testdata/test"

position = 0
res = oss.append_object_from_file(bucket, object, position, localfile)

res = oss.head_object(bucket, object)
if res.status / 100 == 2:
    header_map = convert_header2map(res.getheaders())
    if header_map.has_key('x-oss-next-append-position'):
        position = header_map['x-oss-next-append-position']

res = oss.append_object_from_file(bucket, object, position, localfile)

if res.status / 100 == 2:
    header_map = convert_header2map(res.getheaders())
    if header_map.has_key('x-oss-next-append-position'):
        position = header_map['x-oss-next-append-position']
res = oss.append_object_from_file(bucket, object, position, localfile)
