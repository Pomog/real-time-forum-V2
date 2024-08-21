const newPostElement = (post) => {
    const el = document.createElement("div")
    el.classList.add("post")

    const linkToPost = document.createElement("a")
    linkToPost.classList.add("post-link")
    linkToPost.setAttribute("href", `/post/${post.id}`)
    linkToPost.setAttribute("data-link", "")
    linkToPost.innerText = `${post.title}`

    const postDate = document.createElement("p")
    postDate.innerText = new Date(post.date).toLocaleString()

    const linkToAuthor = document.createElement("a")
    linkToAuthor.setAttribute("href", `/user/${post.author.id}`)
    linkToAuthor.setAttribute("data-link", "")
    linkToAuthor.innerText = `${post.author.firstName} ${post.author.lastName}`

    el.append(linkToPost)
    el.append(postDate)
    el.append(linkToAuthor)

    return el
}

export default {
    newPostElement
}