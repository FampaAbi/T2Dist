# T2Dist

Sebastián Rojas 201773598-8
Fabio Pazos 201773503-1


Para ejecutar:
	- Ejecutar make run en cada maquina ejecutará el código correspondiente.
	-> 	Cliente : make cliente
		  DataNode : make datanode
		  NameNode : make namenode

	- El cliente se puede correr en cualquier máquina asociada a los DATANODES
	- Cabe mencionar, que no porque se haya ejecutado el cliente en una máquina no se puede ejecutar el datanode en la misma, si es que esta es accedida desde otra consola.


Consideraciones:
	-> 	dist61 : datanode
		  dist62 : datanode
		  dist63 : datanode
		  dist64 : namenode

	- Las propuestas a nivel macro se inicializaban al consultar a los nodos si estos se encontraban activos, la respuesta afirmativa de esto produce que esos nodos fuesen considerados. Luego el hecho de aceptar o rechazar iba condicionado a una probabilidad definida por los integrantes (20% rechazo en ambos casos).
  - El apartado de mostrar la disponibilidad de libros, fue considerado de 2 formas ya que no se comprendió bien. Se implementó una opción independiente en el menú para la obtención de estos y de igual manera al acceder al Cliente Downloader desde el menú, son mostrados nuevamente sin una previa solicitud del cliente.
  - En ciertas ocasiones ocurre un error al cual no le pudimos encontrar explicación, las direcciones de las máquinas usadas en la generación de propuestas colapsan y generan un error en cadena al guardarlas de una maner tal como "distdistdist61:9000:9000:9000", en caso de que ocurra esto con un libro, puede realizar 'make reset' y volver a relizar la subida del libro.
