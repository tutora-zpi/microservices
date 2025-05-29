meeting trigger service -> (w przyszlosci) -> meeting scheduler

nie wiem czy tutaj jakiegos cache nie zrobic zeby sprawdzalo te meeting ids'y 

odpalanie eventu za requestowanego z frontu 

spotkanie jest konczone kiedy wszyscy wyjda z niego 

opcje 

- mozna zrobic tutaj prostego schedulera ktory odpalalby spotkanie ktore bylo zaplanowane

- jak bedziemy weryfikowac usera?

- czy podpiac sie do bazy userow i robic read'y 

- gRPC ???, albo tylko weryfikacja podpisu czy jest valid

- mozna dodac baze redis, gdzie sprawdzamy status spotkania 