function download_file() {
    console.log("button is clicked");
    let opts  = {
        method: "get",
        Headers: {
            "Content-Type": "application/json",
        },
        mode: 'cors'
    }

    fetch('http://localhost:8080/generate_download_id', opts).
    then(resp => resp.json()).
    then(obj => {
        console.log("starting download file with download_id = ", obj.download_id);
        let url = "http://localhost:8080/download?download_id=" + obj.download_id;
        window.open(url);
    }).
    catch(err => {
        console.error("error", err)
    });
}