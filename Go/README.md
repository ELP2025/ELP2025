# Projet ELP2025 Go
Nous avons choisi d'implémenter Smoothlife en Go. Smoothlife est une évolution du jeu de la vie. Cette évolution autorise les valeurs flottantes pour l'état des cellules et tente de créer une temporatlité continue.

# Compilation
Il faut télécharger la librairie `fftw` sur son pc avant de compiler

Ensuite, il suffit de faire `go run main` ou `go make main` pour compiler le programme.

# Utilisation 
Options du programme :
- `-i /path/to/image` permet de charger une image comme grille de départ
- `-r` permet de partir d'une grille aléatoire
- `-w et -h` changer la taile de la fenêtre (valeur par défaut 1024x1024)
- `-ra` modifie la taille du kernel. Une valeur plus grande donnera des structures plus grandes.
