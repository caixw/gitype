// 该文件由make.go自动生成，请勿手动修改！

package install

var logFile=[]byte(`<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <info prefix="[INFO]" flag="">
        <console output="stderr" foreground="green" background="black" />
        <buffer size="100">
            <rotate prefix="info-" dir="./output/logs/" size="5M" />
        </buffer>
    </info>

    <debug prefix="[DEBUG]">
        <console output="stderr" foreground="yellow" background="blue" />
        <buffer size="50">
            <rotate prefix="debug-" dir="./output/logs/" size="5M" />
        </buffer>
    </debug>

    <trace prefix="[TRACE]">
        <console output="stderr" foreground="yellow" background="blue" />
        <buffer size="50">
            <rotate prefix="trace-" dir="./output/logs/" size="5M" />
        </buffer>
    </trace>

    <warn prefix="[WARNNING]">
        <console output="stderr" foreground="yellow" background="blue" />
        <rotate prefix="info-" dir="./output/logs/" size="5M" />
    </warn>

    <error prefix="[ERROR]">
        <console output="stderr" foreground="red" background="blue" />
        <rotate prefix="error-" dir="./output/logs/" size="5M" />
    </error>

    <critical prefix="[CRITICAL]">
        <console output="stderr" foreground="red" background="blue" />
        <rotate prefix="critical-" dir="./output/logs/" size="5M" />
    </critical>
</logs>
`)