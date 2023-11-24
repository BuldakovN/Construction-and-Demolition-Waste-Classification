import base64
import os
from flask import Flask, request
#from clientModel import test
import asyncio

app = Flask(__name__)


def treatmenImagetVideo(json):

    #json->video
    fh = open("video/video.mp4", "wb")
    fh.write(base64.b64decode(json["b64"]))
    fh.close()
    answer = {"result":"видео загружено"}

    return answer


@app.route('/sendFromPyVideo', methods=['POST'])
def treatmentVideo():
    return treatmenImagetVideo(request.json)

# @app.route('/start')
# def index():
#     loop = asyncio.new_event_loop()
#     asyncio.set_event_loop(loop)
#     asyncio.get_event_loop().run_until_complete(test())
#     return 'STREAM START'

app.run(host='0.0.0.0', port=5000)