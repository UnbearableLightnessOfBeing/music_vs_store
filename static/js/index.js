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
if (logout) {
  logout.addEventListener("click", async (e) => {
    e.preventDefault()
    try {
      await fetch("/logout", { method: "POST" })
      location.assign(location.origin)
    } catch(err) {
      console.log("error: ", err)
    }
  })
}

// all inputs
const inputsWithNums = document.querySelectorAll("#min_price, #max_price, #label_id")
const inputsWithText = document.querySelectorAll("#price_sorting")
// reset button
const resetButton = document.getElementById("reset")
if (resetButton) {
  resetButton.addEventListener("click", () => {
    if (inputsWithNums && inputsWithNums.length) {
      inputsWithNums.forEach(input => input.value = 0)
    }
    if (inputsWithText && inputsWithText.length) {
      inputsWithText.forEach(input => input.value = "")
    }
  })
}
