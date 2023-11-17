# Harbor Image Tags
> Note:
+ 工具为获取Harbor仓库的某个镜像版本号，结合k8s的kubectl工具发布使用。
+ 代码灵感来源为Jenkins插件项目：[image-tag-parameter-plugin](https://github.com/jenkinsci/image-tag-parameter-plugin)，Jenkins设计k8s的cd可以直接使用此Jenkins插件。
+ 工具可结合Jenkins插件[Active Choices Plug-in](https://plugins.jenkins.io/uno-choice/)使用，但有部分注意点。
## 快速开始
1. 克隆代码&安装
```bash
$ git clone https://github.com/zeratullich/harborgetag
$ cd harborgetag && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build && mv harborgetag /opt
```
2. 使用
```bash
$ /opt/harborgetag -h
Usage of /opt/harborgetag:
  -filter string
        Regular expression to filter image tag e.g. v(\d+\.)*\d+ for tags like v23.3.2 (default ".*")
  -image string
        Harbor registered mirror image (default "qy/book-store")
  -order string
        Allows the user to alter the ordering of the ImageTags in the build parameter.
                Natural-Ordering ... same Ordering as the tags had in prior versions
                Reverse-Natural-Ordering ... the reversed original ordering
                Descending-Versions ... attempts to pars the tags to a version and order them descending
                Ascending-Versions ... attempts to pars the tags to a version and order them ascending
         (default "Descending-Versions")
  -password string
        Harbor auth user password (default "123456")
  -registry string
        Harbor registered mirror address (default "harbor.myquanyi.com")
  -scheme value
        URL request prefix, only one of [http, https] can be selected (default https)
  -username string
        Harbor auth user name (default "admin")
  -verifySSL
        Whether ssl authentication is enabled in harbor (default true)
```
## 与Jenkins插件结合
> 需要首先把编译好的二进制文件放入与Jenkins统一个主机中(vm或者容器)。
1. Jenkins中需要安装[Active Choices Plug-in](https://plugins.jenkins.io/uno-choice/)插件。
![](doc/image/plugin.jpg)
2. 结合参数使用。
![](doc/image/screen1.jpg)
...
3. 需要在最后再应用`Active Choices Plug-in`插件，并需要使用其`Referenced parameters`引入上述自定义的参数。
![](doc/image/screen01.jpg)
其中，如何获取harbor认证账号与密码，请参考此[链接](https://stackoverflow.com/questions/53379151/active-choices-parameter-with-credentials/54927791)。
![](doc/image/screen02.jpg)
4. 点击参数化构建，选择各个可选项，即将可以获取harbor的镜像版本号。
![](doc/image/screen03.jpg)
点击可选项，还可以动态改变获取的镜像标签：
![](doc/image/again01.jpg)
5. 使用Pipeline流水线工作完成CD（持续发布过程）。
示例代码：
```groovy
if(image_tag.contains('error')){
    echo "拉取镜像不正确：${image_tag}"
    error "您选择发布的项目不正确，请正确选择！"
}else{
    def project_name="${project_name}".replaceAll('_', '-')

    echo "您选择发布的环境为：${k8s_cluster}"
    echo "您选择发布的项目为：${project_name}"
    echo "正在发布...."
    if (k8s_cluster == 'DEV'){
        node('qy-xjm-dev-cd'){
            stage('Publish Image'){
                 container('kubectl'){
                    echo "对命名空间：${namespace}的${project_name}项目进行升级或回滚 ！"
                    sh "kubectl set image ${pod_controller} ${project_name} ${project_name}=${registry}/${image_tag} -n ${namespace}"
                    echo "发布完成！" 
                }
            }
        }
    }else if(k8s_cluster == 'UAT'){
        node('qy-xjm-uat-cd'){
            stage('Publish Image'){
                 container('kubectl'){
                    echo "对命名空间：${namespace}的${project_name}项目进行升级或回滚 ！"
                    sh "kubectl set image ${pod_controller} ${project_name} ${project_name}=${registry}/${image_tag} -n ${namespace}"
                    echo "发布完成！" 
                }
            }
        }
    }else if(k8s_cluster == 'PROD'){
        node('qy-xjm-prod-cd'){
            stage('Publish Image'){
                 container('kubectl'){
                    echo "对命名空间：${namespace}的${project_name}项目进行升级或回滚 ！"
                    sh "kubectl set image ${pod_controller} ${project_name} ${project_name}=${registry}/${image_tag} -n ${namespace}"
                    echo "发布完成！" 
                }
            }
        }
    }
}
```
其中上述代码中的`qy-xjm-dev-cd`、`qy-xjm-uat-cd`、`qy-xjm-prod-cd`等`node`信息需要在Jenkins中使用kubernetes插件设置Cloud中定义，确保jenkins能连接到各个k8s集群(如何创建请咨询)。
![](doc/image/screen04.jpg)
...

**注意:** Jenkins如果可以安装[image-tag-parameter-plugin](https://github.com/jenkinsci/image-tag-parameter-plugin)插件（需要Jenkins版本比较高），可使用此插件，因为出现错误时候可以参照Jenkins日志查看错误原因，使用[Active Choices Plug-in](https://plugins.jenkins.io/uno-choice/)插件无法给出错误日志，故排查问题比较困难，但此插件灵活性差，不便于项目的联动发布。此工程可以用在其他需要获取Harbor仓库镜像版本号的应用上，例如使用脚本获取镜像版本号来作为持续发布（CD）的方案上。如有实际CICD帮助需要(有整套解决方案，已总结成文档)请联系：lichuan4961(微信)，或直接提issue，加微信请标注来意。
## 手动执行获取tag示例
![](doc/image/screen05.jpg)