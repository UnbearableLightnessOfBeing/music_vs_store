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

// product page
let price = 0
const priceEl = document.getElementById("product-price")
let priceStr = ""
if (priceEl) {
  priceStr = priceEl.innerText
}
const split = priceStr.split(" ")
if (split.length) {
  const decimal = split[0]
  const decSplit = decimal.split(".")
  if (decSplit.length) {
    price = decSplit[0]
  }
}

const totalEl = document.getElementById("total")

const setTotalPrice = (quantity) => {
  if (totalEl) {
    totalEl.innerText = `${Number(quantity) * price}.00 руб`
  }
}

const quantity = document.getElementById("quantity")
const inc = document.getElementById("increment-quantity")
const dec = document.getElementById("decrement-quantity")
if (quantity && inc && dec) {
  inc.addEventListener("click", () => {
    quantity.value++
    setTotalPrice(quantity.value)
  })

  dec.addEventListener("click", () => {
    if (quantity.value > 1) {
      quantity.value--
      setTotalPrice(quantity.value)
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

const mainCarousel = document.getElementById("main-carousel")
const thumbnails = document.getElementById("thumbnail-carousel")

if (mainCarousel && thumbnails) {
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

  main.sync(thumnails)
  main.mount()
  thumnails.mount()
}

const brands = document.getElementById("brands")
if (brands) {
  const splide = new Splide( '#brands', {
    type   : 'loop',
    drag   : 'free',
    focus  : 'center',
    perPage: 4,
    arrows: false,
    pagination: false,
    autoScroll: {
      speed: 1,
    },
  } );

  splide.mount(window.splide.Extensions);
}


