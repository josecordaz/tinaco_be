from flask import Flask
import RPi.GPIO as GPIO
import time

GPIO.setmode(GPIO.BCM)

TRIG = 23
ECHO = 24

app = Flask(__name__)

@app.route("/")
def hello():
    return str(getMeasure())


def getMeasure():
    return 45

if __name__ == '__main__':
    app.run(host= '0.0.0.0',debug=True)