from flask import Flask
import RPi.GPIO as GPIO
import time
import random
import json

GPIO.setmode(GPIO.BCM)

TRIG = 23
ECHO = 24

app = Flask(__name__)

data = {
    "president": {
        "name": "Zaphod Beeblebrox",
        "species": "Betelgeusian"
    }
}

@app.route("/")
def hello():
    return json.dumps(data)

@app.route("/level")
def level():
    dt = {
        "level":str(getMeasure())
    }
    return json.dumps(dt)

@app.route("/b_on")
def level():
    return "encender bomba"

@app.route("/b_off")
def level():
    return "apagar bomba"

@app.route("/b_status")
def level():
    return "bomba estatus"

# mux.HandleFunc("/login", login)
# mux.HandleFunc("/level", getlevel)
# mux.HandleFunc("/b_on", TurnBombOn)
# mux.HandleFunc("/b_off", TurnBombOff)
# mux.HandleFunc("/b_status", GetBombStatus)

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