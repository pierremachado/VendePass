import json

# Função para ler um arquivo JSON
def read_json(filename):
    with open(filename, 'r') as file:
        return json.load(file)

# Função para criar o grafo de adjacências
def create_graph(flights):
    graph = {}
    for flight in flights:
        src = flight['SourceAirportId']
        dest = flight['DestAirportId']
        
        if src not in graph:
            graph[src] = []
        if dest not in graph:
            graph[dest] = []
        
        graph[src].append(dest)
    return graph

# Função para imprimir as adjacências do grafo
def print_adjacencies(graph, airports):
    airport_map = {airport['Id']: airport['City']['Name'] for airport in airports}
    
    for airport_id, destinations in graph.items():
        city_name = airport_map.get(airport_id, 'Unknown')
        dest_cities = [airport_map.get(dest_id, 'Unknown') for dest_id in destinations]
        print(f"{city_name} is connected to: {', '.join(dest_cities)}")

def main():
    # Lê os arquivos JSON
    flights = read_json('internal/stubs/flights.json')
    airports = read_json('internal/stubs/airports.json')
    
    # Cria o grafo
    graph = create_graph(flights)
    
    # Imprime as adjacências
    print_adjacencies(graph, airports)

if __name__ == "__main__":
    main()
