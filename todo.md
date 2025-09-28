- [x] auto update dans le bin
- [x] ajouter des tests ( voir quoi et comment tester)
      => et dans la ci

- [x] ci : voir tout ce qu'on peut mettre qui a du sens.

- [ ] autocomplete

- [ ] faire un truc sympa avec les updates.
      => il faudrait ajouter un changelog, et quand on update, il detecte le diff entre la version installée et la nouvelle, et présente chacun des blocs ( avec des actions possibles à chaque fois ? notamment autour de la config qui pourrait être outdated ou incompatible)

- [x] pouvoir gérer un debug mode pour printer les logs ( via la config, ou via process.env equivalent )
- [x] ajouter une commande pour créer le fichier de conf local et pousser le path dans le gitignore

par contre y'a plusieurs trucs qui vont pas : tes exemples avec echo sont pas bon, c'est
pas comme ça que ça marchera dans le yml.\
 ensuite, il faut referencer hexa dans la commande, pas la version locale.\
 ensuite pour la partie setup de config, c'est un truc qu'on va faire là, mais du coup la
doc devrait le notifier :\
 le post install / update via homebrew devra lancer la commande qui copie le fichier de
template de conf au niveau user.\
 \
 la doc devrait en parler, on va attaquer ce morceau après.
