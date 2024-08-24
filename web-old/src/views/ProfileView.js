import fetcher from "../services/Fetcher.js";
import AbstractView from "./AbstractView.js";
import PostUtils from "../services/PostUtils.js"

const genders = { 1: 'Male', 2: 'Female' }
const roles = {  1: 'Guest',  2: 'User', 3: 'Moderator',  4: 'Administrator' }

const getUserByID = async (id) => {
    const path = `/api/users/${id}`
    return await fetcher.get(path);
}

const getUsersPosts = async (userID) => {
    const path = `/api/users/${userID}/posts`
    return await fetcher.get(path);
}

const getUsersRatedPosts = async (userID) => {
    const path = `/api/users/${userID}/rated-posts`
    return await fetcher.get(path);
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Profile");
        this.userID = params.userID;
    }

    async getHtml() {
        return `
        <h2>Users profile</h2>
        <div id="user-profile">
                <div class="profile-info" id="avatar"></div>
                <div>
                    <div class="profile-info" id="username"></div>
                    <div class="profile-info" id="first-name"></div>
                    <div class="profile-info" id="last-name"></div>
                    <div class="profile-info" id="age"></div>
                    <div class="profile-info" id="gender"></div>
                    <div class="profile-info" id="role"></div>
                    <div class="profile-info" id="registered"></div>
                </div>
            </div>
            <h2>Users posts</h2>
            <div id="users-posts"></div>
            <h2>Users liked posts</h2>
            <div id="users-liked-posts"></div>
            <h2>Users disliked posts</h2>
            <div id="users-disliked-posts"></div>
        `;
    }

    async init() {
        const user = await getUserByID(this.userID)
       
        document.querySelector('.profile-info#avatar').innerHTML = `<img src="http://${API_HOST_NAME}/images/${user.avatar}">`
        document.querySelector('.profile-info#username').innerText = `Username: ${user.username}`
        document.querySelector('.profile-info#first-name').innerText = `First name: ${user.firstName}`
        document.querySelector('.profile-info#last-name').innerText = `Last name: ${user.lastName}`
        document.querySelector('.profile-info#age').innerText = `Age: ${user.age}`
        document.querySelector('.profile-info#gender').innerText = `Gender: ${genders[user.gender]}`
        document.querySelector('.profile-info#role').innerText = `Role: ${roles[user.role]}`
        document.querySelector('.profile-info#registered').innerText = `Registered: ${new Date(Date.parse(user.registered)).toLocaleString()}`

        const usersPosts = await getUsersPosts(this.userID) 
        const usersRatedPosts = await getUsersRatedPosts(this.userID)|| []

        const usersPostsEl = document.getElementById('users-posts')
        if (usersPosts != null) {
            usersPosts.forEach((post) => {
                const postEl = PostUtils.newPostElement(post);
                usersPostsEl.append(postEl)
            })
        } else {
            usersPostsEl.innerText = 'No posts'
        }


        const usersLikedPosts = usersRatedPosts.filter((post) => post.userRate == 1)
        const usersLikedPostsEl = document.getElementById('users-liked-posts')
        if (usersLikedPosts.length > 0 ) {
            usersLikedPosts.forEach((post) => {
                const postEl = PostUtils.newPostElement(post);
                usersLikedPostsEl.append(postEl)
            })
        } else {
            usersLikedPostsEl.innerText = 'No posts'
        }

        const usersDisLikedPosts = usersRatedPosts.filter((post) => post.userRate == 2)
        const usersDislikedPostsEl = document.getElementById('users-disliked-posts')
        if (usersDisLikedPosts.length > 0) {
            usersDisLikedPosts.forEach((post) => {
                const postEl = PostUtils.newPostElement(post);
                usersDislikedPostsEl.append(postEl)
            })
        } else {
            usersDislikedPostsEl.innerText = 'No posts'
        }
    }
}