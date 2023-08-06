function search(){
    let search = document.getElementById('searchbar').value.toUppercase();
    let product = document.getElementById('pro');
    let pro = document.querySelector('section-1');
    const name = document.getElementsByTagName('h5');


    for(var i=0;i<name.length;i++){
        let match = pro[i].getElementsByTagName('h5')[0];

if(match){
    let textvalue = match.textContent || match.innerHTML
    

    if(textvalue.toUppercase().indexOf(search)>-1){
        pro[i].style.display="";

    }else{
        pro[i].style.display="none";
    }
}

    }
}