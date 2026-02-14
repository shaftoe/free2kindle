<script lang="ts">
	import { resolve } from '$app/paths';
</script>

<nav>
	<ul>
		<li><a href={resolve('/')}>My List</a></li>
		<li><a href={resolve('/new')}>Save new article</a></li>
		<li><a href={resolve('/settings')}>Settings</a></li>
	</ul>
</nav>
