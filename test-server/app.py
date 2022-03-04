
from flask import Flask

app = Flask(__name__)
@app.route('/add-event', methods=['POST'])
def create_event():
    pass


if __name__ == '__main__':
    app.run(debug=True)