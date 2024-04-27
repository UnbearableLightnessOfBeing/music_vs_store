// menu switching
const menu = document.getElementById("menu-dropdown")
const userIcon = document.getElementById("user-icon")

userIcon.addEventListener("click", function(e) {
  e.stopPropagation()
  menu.classList.remove("hidden")
  this.classList.add('text-amber-600')
})

menu.addEventListener("click", (e) => e.stopPropagation())

window.addEventListener("click", () => {
  menu.classList.add("hidden")
  userIcon.classList.remove('text-amber-600')
})

// logout link
const logout = document.getElementById("logout")
logout.addEventListener("click", async (e) => {
  e.preventDefault()
  try {
    await fetch("/logout", { method: "POST" })
    location.assign(location.origin)
  } catch(err) {
    console.log("error: ", err)
  }
})
