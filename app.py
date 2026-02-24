from flask import Flask, request, jsonify
from datetime import datetime

app = Flask(__name__)

users = []

@app.route('/usuarios', methods=['POST'])
def create_user():
    data = request.get_json()

    if not data or 'nome' not in data or 'sobrenome' not in data or 'email' not in data or 'dataNascimento' not in data:
        return jsonify({'message': 'Missing required fields'}), 422

    nome = data['nome']
    sobrenome = data['sobrenome']
    email = data['email']
    dataNascimento = data['dataNascimento']

    if not isinstance(nome, str) or not (2 <= len(nome) <= 50) or not all(x.isalpha() or x.isspace() for x in nome):
        return jsonify({'message': 'Invalid nome'}), 422

    if not isinstance(sobrenome, str) or not (2 <= len(sobrenome) <= 50) or not all(x.isalpha() or x.isspace() for x in sobrenome):
        return jsonify({'message': 'Invalid sobrenome'}), 422

    if not isinstance(email, str) or '@' not in email:
        return jsonify({'message': 'Invalid email format'}), 422

    for user in users:
        if user['email'] == email:
            return jsonify({'message': 'Email already exists'}), 400

    try:
        birthdate = datetime.strptime(dataNascimento, '%Y-%m-%d').date()
    except ValueError:
        return jsonify({'message': 'Invalid dataNascimento format. Use YYYY-MM-DD'}), 422

    today = datetime.now().date()
    age = today.year - birthdate.year - ((today.month, today.day) < (birthdate.month, birthdate.day))

    if birthdate >= today:
        return jsonify({'message': 'Invalid dataNascimento. Must be in the past.'}), 422

    if age < 18:
        return jsonify({'message': 'User must be at least 18 years old.'}), 422

    try:
        user_id = len(users) + 1
        user = {
            'id': user_id,
            'nome': nome,
            'sobrenome': sobrenome,
            'email': email,
            'dataNascimento': dataNascimento
        }
        users.append(user)
        return jsonify({'id': user_id, 'user': user}), 200
    except Exception as e:
        return jsonify({'message': str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True, port=5000)
