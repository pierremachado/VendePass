import socket
import json

# Configurações
ADDRESS = 'localhost'
PORT = 8080

def send_request(action, data):
    """Envia uma solicitação ao servidor e retorna a resposta."""
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        try:
            # Conecta ao servidor
            sock.connect((ADDRESS, PORT))
            
            # Cria o pedido
            request = {
                "Action": action,
                "Data": data
            }
            
            # Serializa o pedido em JSON
            buffer = json.dumps(request).encode('utf-8')
            
            # Envia o pedido ao servidor
            sock.sendall(buffer)
            
            # Recebe a resposta do servidor
            receive = sock.recv(2048)
            
            # Deserializa a resposta JSON
            response = json.loads(receive.decode('utf-8'))
            return response
            
        except socket.error as e:
            print(f"Erro na conexão ou na comunicação: {e}")
            return None

def main():
    # Cria a conexão para o login
    login_response = send_request(
        "login",
        {
            "Username": "pedrocosta",
            "Password": "senhaSegura79"
        }
    )
    
    if login_response:
        print("Resposta de Login:", login_response)
        token = login_response.get("Data", {}).get("token", "")
        
        # Cria a conexão para o logout
        if token:
            logout_response = send_request(
                "logout",
                {
                    "TokenId": token
                }
            )
            print("Resposta de Logout:", logout_response)
        else:
            print("Token não encontrado para o logout")
    else:
        print("Erro ao receber a resposta de login")

if __name__ == '__main__':
    main()
