// menu switching
const menu = document.getElementById("menu-dropdown")
const userIcon = document.getElementById("user-icon")

userIcon.addEventListener("click", function(e) {
  e.stopPropagation()
  menu.classList.remove("hidden")
  this.classList.add('text-secondary')
})

menu.addEventListener("click", (e) => e.stopPropagation())

window.addEventListener("click", () => {
  menu.classList.add("hidden")
  userIcon.classList.remove('text-secondary')
})

// logout link
const logout = document.getElementById("logout")
if (logout) {
  logout.addEventListener("click", async (e) => {
    e.preventDefault()
    try {
      await fetch("/logout", { method: "POST" })
      location.assign(location.origin)
    } catch (err) {
      console.log("error: ", err)
    }
  })
}

// all inputs
const inputsWithNums = document.querySelectorAll("#label_id")
const inputsWithText = document.querySelectorAll("#price_sorting, #min_price, #max_price")
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

// splide
// new Splide(".splide", {
//   perPage: 4,
//   perMove: 1,
//   gap: 10,
//   pagination: false,
//   arrows: false,
// }).mount();
//

const main = new Splide("#main-carousel", {
  type: "slide",
  rewind: true,
  pagination: false,
  arrows: true,
  perMove: 1,
})

const thumnails = new Splide('#thumbnail-carousel', {
  fixedWidth: 100,
  fixedHeight: 80,
  perPage: 4,
  gap: 10,
  rewind: true,
  pagination: false,
  arrows: false,
  focus: 'center',
  isNavigation: true,
  breakpoints: {
    600: {
      fixedWidth: 60,
      fixedHeight: 44
    }
  }
});

if (main && thumnails) {
  main.sync(thumnails)
  main.mount()
  thumnails.mount()
}

