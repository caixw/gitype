<?xml version="1.0" encoding="utf-8"?>
<!--
为sitemap.xml产生一个比较美观的人机界面。

@author     caixw <http://github.com/caixw>
@copyright  Copyright(C) 2010-2015, caixw
@license    MIT License
@date       2010-01-02
@update     2015-10-20
-->
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" xmlns:sm="http://www.sitemaps.org/schemas/sitemap/0.9">
<xsl:output method="html" encoding="utf-8" indent="yes" version="1.0" />
<xsl:template match="/">
<!-- xsl:text disable-output-escaping='yes'>&lt;!DOCTYPE html&gt;</xsl:text -->
<html>
<head>
<title>XML Sitemap</title>
<meta charset="utf-8" />
<meta name="generator" content="http://caixw.io" />
<style type="text/css">
a{text-decoration:none;color:#123}
a:hover{text-decoration:underline;color:#c96}
.bold{font-weight:bold}
.error{color:red}

header h1{font-size:1.5em;font-weight:bold}
header .desc,footer{margin:0.7em;line-height:1.8em}
header a,footer a{color:blue}

table{width:100%;text-align:left;line-height:1.5em;border-collapse:collapse}
td, th{padding:0em 0.3em}
thead tr, tfoot tr{background:#ddd;height:1.6em}
tbody tr:nth-of-type(even){background:#eee}
tbody tr:hover{background:#ddd}
</style>
</head>
<body>
    <header>
    <h1>XML Sitemap</h1>
    <div class="desc">这是个标准的sitemap文件。您可以将此文件提交给<a target="_blank" href="http://www.google.com/webmasters/tools/">Google</a>、<a target="_blank" href="http://www.bing.com/webmaster">Bing</a>或<a target="_blank" href="http://sitemap.baidu.com">百度</a>，让搜索引擎更好地收录您的网站内容。<br />
        若是存在sitemap的索引文件，则<span class="bold">只需提交索引文件</span>即可。更详细的信息请参考<a href="http://www.sitemaps.org/zh_CN/protocol.php">这里</a>。
        </div><!-- end desc -->
    </header>
    <xsl:apply-templates select="sm:urlset" />
    <footer>此XSL模板由<a target="_blank" href="https://github.com/caixw">caixw</a>制作，并基于<a target="_blank" href="http://www.opensource.org/licenses/MIT">MIT</a>版权发布。</footer>
</body>
</html>
</xsl:template>


<xsl:template match="sm:urlset">
<div id="content">
<table>
    <thead>
    <tr>
        <th>地址</th>
        <th>最后更新</th>
        <th>更新频率</th>
        <th>权重</th>
    </tr>
    </thead>
    <tfoot>
        <tr><td colspan="4">当前总共<xsl:value-of select="count(/sm:urlset/sm:url)" />条记录</td></tr>
    </tfoot>
    <tbody>
        <xsl:for-each select="sm:url">
        <tr>
            <td><a>
                <xsl:attribute name="href"><xsl:value-of select="sm:loc" /></xsl:attribute>
                <xsl:value-of select="sm:loc" />
            </a></td>
            <td><xsl:value-of select="concat(substring-before(sm:lastmod, 'T'),' ',substring(sm:lastmod,12,5))" /></td>
            <td>
                <xsl:choose>
                    <xsl:when test="sm:changefreq = 'never'">从不</xsl:when>
                    <xsl:when test="sm:changefreq = 'yearly'">每年</xsl:when>
                    <xsl:when test="sm:changefreq = 'monthly'">每月</xsl:when>
                    <xsl:when test="sm:changefreq = 'weekly'">每周</xsl:when>
                    <xsl:when test="sm:changefreq = 'daily'">每天</xsl:when>
                    <xsl:when test="sm:changefreq = 'hourly'">每小时</xsl:when>
                    <xsl:when test="sm:changefreq = 'always'">实时</xsl:when>
                    <xsl:otherwise><span class="error">未知的值</span></xsl:otherwise>
                </xsl:choose>
            </td>
            <td><xsl:value-of select="concat(sm:priority*100,'%')" /></td>
        </tr>
        </xsl:for-each>
    </tbody>
</table>
</div>
</xsl:template>

</xsl:stylesheet>
