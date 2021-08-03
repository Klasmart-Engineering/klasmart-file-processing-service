# kidsloop-file-processing-service
In kidsloop2, students always upload their files into S3 as homework, which maybe lead to the leakage of sensitive personal information of students. In order to protect the security of our users in kidsloop2, we need to hide some metadata of the files in S3 such as geographic information in students' profile, computer names in the pdf files and etc.
To protect users' sensitive infomation, we can postpone hide the metadata when the users upload the files. For example when a student has uploaded his profile, his profile hasn't changed yet. He should click a save button to save his change. We can do the hide work when the student save his change.


![image.png](https://cdn.nlark.com/yuque/0/2021/png/312446/1627969728514-e6f4b268-7f93-4d0e-ad9f-2e5fcf54ffa9.png)

The process of the solution is when after the user click the upload file, front-end program send the request for the signed url for S3 and upload the file via the signed url. Currently the file is S3 contains the sensitive metadata. However the uploaded file hasn't record in the database, we can't access it via Kidsloop2. We add some code in the save process in the back-end service to hide the sensitive metadata. When the user click the save button, we save the change while access S3 for the file, hide the sensitive metadata and overwrite it into S3.
The file processing will take a lot of time. So we deploy a brand new service "**kidsloop-file-processing-service**" to do the job.
