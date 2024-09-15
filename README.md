## guiEDF

### But

Cette interface graphique permet d'analyser la consommation EDF sur une année.

Ceci afin de comparer les options de Base,Heures Creuses et Tempo.

### Prérequis

*   Sur le side Enedis  [https://www.enedis.fr/](https://www.enedis.fr/)  
    Créer un compte si nécessaire
*   Se connecter  
       Suivre mes mesures / Télécharger mes données :  
           Type de données **: Consommation horaire** (ceci est   
           nécessaire pour analyser les heures creuses et Tempo)  
       Période choisie:   
           Sélectionner une date de début et une date de fin pour   
           une période supérieure à une année.  
    Télécharger mes données

### _**Enedis enverra ultérieurement un fichier de type .csv**_

Application

*   **Fich Enedis.cvs:**   
       Sélectionner le fichier Enedis prédemment téléchargé.  
       Les dates de Début et de Fin Enedis seront mises à jour.
*   **Date Début et Date de Fin:**   
       Entrer une heure de fin incluse dans la période du fichier   
       Enedis.  
       La date de Début sera mise à jour automatiquemen (Fin - 1 an)  
       Choisir impérativement une date de Fin pour avoir une période   
       d'analyse connue du fichier Enedis.
*   Les jours Tempo de la période d'analyse seront éventuellement   
       mis à jour automatiquement via le site EDF.
*   **Puissance Installée:**   
       Choisir la puissance installée de votre compteur  
       Ceci influera sur les tarifs (abonnement et consommation)
*   **Créneaux H.Creuses:**   
       Définir les créneaux Heures Creuses  
       (3 créneaux max) le format doit être:  
       hh:mn-hh:mn  
           hh:mn-hh:mn ou   
           hh:mn-hh:mn,hh:mn-hh:mn ou   
           hh:mn-hh:mn,hh:mn-hh:mn,hh:mn-hh:mn  
       Le bon formatage de votre saisie sera indiqué
*   **Tarifs EDF au:**    
       Affichage de la date des derniers tarifs connus  
       M. à jour:   
           Les derniers tarifs seront récupérés chez EDF et la date   
           sera mise à jour
*   **Calculer:**   
       Analyse d'une année d'enregistrement EDF  
       Une alerte sera générée si les créneaux Heures Creuses ne sont   
       pas connus ou incorrects.  
       Le bouton défaut sélectionnera un créneau par défaut 22:00-0600

**Les résultats seront affichés**

*   **Quit:**   
       Quitter l'application.
*   **Aide:**   
       Afficher ce fichier d'aide..