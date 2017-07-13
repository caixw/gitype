// 这是自动产生的文件，不需要修改

package admin

var AdminHTML = `<!DOCTYPE html>
<html lang="zh-cmn-Hans">
	<meta charset="utf-8" />
	<title>typing 控制面板</title>
    <style>
        body{text-align:center}
        .container{
            margin:auto;
            margin-top:5rem;
            text-align:left;
            width:30rem;
        }

        form input,form button{font-size:1.2rem}

        a{text-decoration:none}
    </style>
	<body>
	<div class="container">
		<h1>控制面板</h1>
		<p>
			<span>最后更新时间:</span>{{.lastUpdate}}
		</p>

		<form action="" method="POST">
			<p>
				<input type="password" name="password" placeholder="密码" />
				<button type="submit">重新生成</button>
			</p>
		</form>
	</div>
	</body>
</html>
`