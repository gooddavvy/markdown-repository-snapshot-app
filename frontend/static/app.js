document.addEventListener("DOMContentLoaded", function () {
    let ignoreList = [];
    let randomChoice = function (array) {
        return array[Math.floor(Math.random() * array.length)];
    };
    let isValidURL = function (url) {
        // Regular expression for URL validation
        var urlPattern = /^(https?:\/\/)?([\w-]+\.)+[\w-]+(\/[\w-./?%&=]*)?$/;
        return urlPattern.test(url);
    };
    let showModal = function (content) {
        let modal = document.getElementById("modal");
        let modalContent = document.getElementById("modal-content");
        modalContent.innerHTML = content;
        modal.classList.remove("hidden");
        modal.style.display = "block";

        // Close modal when clicking on close button or outside the modal
        let closeModal = document.querySelector(".close-modal");
        modal.onclick = function (event) {
            if (event.target === modal || event.target === closeModal) {
                hideModal();
            }
        };
        closeModal.onclick = () => hideModal()
    };
    let hideModal = function () {
        let modal = document.getElementById("modal");

        modal.style.removeProperty("display");
        modal.classList.add("hidden");
    };
    let createModal = function (innerHTML) {
        showModal(innerHTML);
    };



    document.getElementById("ignore_item").placeholder = randomChoice([".github", "README.md", "remote"]);

    document.getElementById("add_ignore_item").addEventListener("click", function () {
        let ignoreItem = document.getElementById("ignore_item").value.trim();
        let listLabel = document.getElementById("list_label");
        if (ignoreItem) {
            listLabel.classList.remove("hidden");
            ignoreList.push(ignoreItem);
            let li = document.createElement("li");
            li.textContent = ignoreItem;
            document.getElementById("ignore_list").appendChild(li);
            document.getElementById("ignore_item").value = "";
        }
    });

    document.getElementById("generate_markdown").addEventListener("click", function () {
        let generateButton = document.getElementById("generate_markdown");
        let repositoryUrl = document.getElementById("repository_url").value.trim();

        if (!repositoryUrl || !isValidURL(repositoryUrl)) {
            createModal(`
            <h2>Missing or invalid repository URL</h2>
            <p>A valid repository URL is required.</p>
            `);
            return;
        }

        generateButton.disabled = true;
        generateButton.classList.add("disabled");
        generateButton.textContent = "Generating, please wait...";

        fetch("/generate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ repository_url: repositoryUrl, ignore_list: ignoreList })
        })
            .then(response => response.json())
            .then(data => {
                let responseSection = document.getElementById("response_section");
                let downloadLink = document.getElementById("download_link");
                let responseJson = document.getElementById("response_json");

                // Ensure response content and download link are hidden initially
                responseJson.classList.add("hidden");
                downloadLink.classList.add("hidden");

                // Show the response section
                responseSection.classList.remove("hidden");

                // Handle download link visibility
                if (typeof data.response === "string") {
                    let blob = new Blob([data.response], { type: "text/markdown" });
                    let url = URL.createObjectURL(blob);

                    downloadLink.href = url;
                    downloadLink.classList.remove("hidden");

                    generateButton.textContent = "Finished generating, thank you for your patience.";
                } else {
                    responseJson.textContent = JSON.stringify(data.response, null, 2);
                    responseJson.classList.remove("hidden");
                    generateButton.textContent = "An error occurred while generating, sorry about that.";
                }
            });
    });
});