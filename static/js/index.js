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

// filtering and sorting
const sortingParams = {
  min_price: 0,
  max_price: 0,
  label_id: 0
}

const replaceRouteWithQueries = () => {
  const result = Object.keys(sortingParams).reduce((acc, cv, idx) => {
    if (sortingParams[cv] > 0) {
      acc += (idx !== 0 ? "&" : "" + cv + "=" + sortingParams[cv].toString())
      return acc
    } 
    return acc
  }, "")

  const base = location.href.split("?")[0]
  location.replace(base + "?" + result)
}

// min-price input
const minPriceInput = document.getElementById("min-price")
if (minPriceInput) {
  minPriceInput.addEventListener("change", (e) => {
    sortingParams.min_price = Number(e.target.value)
  })
}

// apply button
const applyButton = document.getElementById("apply")
if (applyButton) {
  applyButton.addEventListener("click", replaceRouteWithQueries)
} 

