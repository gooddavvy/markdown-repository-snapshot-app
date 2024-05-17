document.addEventListener("DOMContentLoaded", function () {
    let ignoreList = [];
    let randomChoice = function (array) {
        return array[Math.floor(Math.random() * array.length)];
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
        document.getElementById("generate_markdown").disabled = true;
        document.getElementById("generate_markdown").classList.add("disabled");
        document.getElementById("generate_markdown").textContent = "Generating, please wait...";

        let repositoryUrl = document.getElementById("repository_url").value;
        if (repositoryUrl) {
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
                    responseSection.classList.add("hidden");
                    downloadLink.classList.add("hidden");

                    // Show the response section
                    responseSection.classList.remove("hidden");

                    // Handle download link visibility
                    if (typeof data.response === "string") {
                        let blob = new Blob([data.response], { type: "text/markdown" });
                        let url = URL.createObjectURL(blob);

                        downloadLink.href = url;
                        downloadLink.classList.remove("hidden");

                        document.getElementById("generate_markdown").textContent = "Finished generating, thank you for your patience.";
                    } else {
                        downloadLink.classList.add("hidden");
                        responseJson.classList.remove("hidden");

                        responseJson.textContent = JSON.stringify(data.response);
                        document.getElementById("generate_markdown").textContent = "An error occurred while generating, sorry about that.";
                    }
                });
        }

        // document.getElementById("generate_markdown").textContent = "Finished generating, thank you for your patience.";
    });
});
