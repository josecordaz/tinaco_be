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
    return random.random() * 100
    # print "Distance Measurement In Progress"

    # GPIO.setup(TRIG,GPIO.OUT)
    # GPIO.output(TRIG,0)

    # GPIO.setup(ECHO,GPIO.IN)

    # time.sleep(0.1)

    # print "Stargin measurement"

    # GPIO.output(TRIG, 1)
    # time.sleep(0.00001)
    # GPIO.output(TRIG, 0)

    # while GPIO.input(ECHO) == 0:
    #     pass
    # start = time.time()

    # while GPIO.input(ECHO) == 1:
    #     pass
    # stop = time.time()

    # print (stop - start)

    # print (stop - start) * 17000

    # GPIO.cleanup()

    # return (stop - start) * 17000

if __name__ == '__main__':
    app.run(host= '0.0.0.0',debug=True)