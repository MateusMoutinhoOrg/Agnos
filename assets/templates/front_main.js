// Every script under assets/frontend/static/scripts is linked by `dirref
// "scripts"`, sorted by path and deferred, so this file runs after the
// document has parsed. It is yours from here: front-init writes it once and
// never reads it again.
console.log("front layer ready");
