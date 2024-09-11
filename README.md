# -
本项目是我在系统性的学习了GO语言的基础之后所作的一个gin-grom项目，在这个项目之中的前端的代码全部是用ChatGPT进行生成的，导致了前端页面非常的不美观，希望有全栈的大佬能对该项目的前端做出优化。
在测试接口的适合建议使用postman
本项目的所有代码放在了app这个包里面，view文件夹下面放的是部分的前端的简单的代码：
![image](https://github.com/user-attachments/assets/0123ab32-5ebd-40db-90f1-430ad99ec2ed)
后端的逻辑接口全部放在了logic里面，里面包含了管理员对图书的增删改查等操作等等：
![image](https://github.com/user-attachments/assets/d783228f-dc2c-44df-9665-9be298a93c6e)
然后执行逻辑的数据库接口全部放在了model文件夹的下面：
![image](https://github.com/user-attachments/assets/0f90501b-d463-4652-a96a-f21c4fc52795)
在这个项目之中，我们运用到了大量的基础知识，我现在列举出主要的知识点：
1.设置session保存用户的登录态。
2.如何实现验证码登陆  1.数字验证码和邮箱验证码以及手机验证码。 
3.实现了一个限流策略 XYZ限流
4.在MySQL中存储用户密码的适合，是通过md5加盐实现的。
5.使用雪花算法生成唯一的uid
6.通过沙盒来模拟支付宝支付
