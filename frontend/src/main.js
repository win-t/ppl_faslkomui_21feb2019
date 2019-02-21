let counter = document.getElementById('counter');
let username = document.getElementById('username');
let password = document.getElementById('password');
let reset = document.getElementById('reset');

reset.onclick = () => {
  (async () => {
    let res = await fetch(API_URL + '/reset', {
      method: 'POST',
      headers: (()=>{
        let h = new Headers();
        h.set('Authorization', 'Basic ' + btoa(username.value + ":" + password.value));
        return h;
      })(),
    });
    if (res.status != 200) {
      throw new Error(`API not returning 200, but ${res.status}\n${await res.text()}`);
    }
    counter.innerText = '0';
  })().catch(err => {
    console.error(err);
    alert(`ERR: ${err.message}`);
  })
}

(async () => {
  let res = await fetch(API_URL + '/counter');
  if (res.status != 200) {
    throw new Error(`API not returning 200, but ${res.status}`);
  }
  res = await res.json();
  counter.innerText = res.counter;
})().catch(err => {
  console.error(err);
  alert(`ERR: ${err.message}`);
});
