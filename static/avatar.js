const recherche = document.getElementById("rechercheAvatar");
const avatars = document.querySelectorAll(".avatar-choix");

recherche.addEventListener("input", function () {

    const texte = this.value.toLowerCase();

    avatars.forEach(function (avatar) {

        const nom = avatar.dataset.nom.toLowerCase();

        if (nom.startsWith(texte)) {
            avatar.style.display = "";
        } else {
            avatar.style.display = "none";
        }

    });

});

avatars.forEach(function (avatar) {
    const choix = avatar.querySelector('input[name="avatar"]');
    choix.addEventListener("change", function () {
        if (this.checked) {
            document.body.style.setProperty("--character-background", 'url("' + avatar.dataset.background + '")');
            this.form.requestSubmit();
        }
    });
});
