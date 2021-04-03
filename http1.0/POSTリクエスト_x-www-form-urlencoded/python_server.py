from flask import Flask, render_template, request

app = Flask(__name__)


@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        dump_request()

    return render_template('index.html')


def dump_request():
    content_type = request.headers.get('Content-Type')
    request_body = request.get_data()

    print(f'POST {request.path} {request.headers.environ.get("SERVER_PROTOCOL")}')
    print(f'Host: {request.host}')
    print(f'Content-Type: {content_type}')
    print()
    print(request_body.decode('utf-8'))


def main():
    app.run(debug=True)


if __name__ == '__main__':
    main()
