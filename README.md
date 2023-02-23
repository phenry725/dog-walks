# Dog Walking Calculator



We have a dog walker who isn't the most digital of a person. So, we kept having to go back and check out which days the walks had happened and we would frequently pay the dog walker incorrectly.

Now this program takes calendar events that have a certain prefix as a title ![Event Example](https://github.com/phenry725/dog-walks/blob/main/images/dogWalk.JPG?raw=true)

(default is "[DOG WALK]" but that is configurable) and counts them and applies the rate. Now, as long as you make sure the walks are accurate on the calendar, this should make sure the dog walker is paid correctly.

In order to avoid having to constantly fetch a new OAuth token to pull from your calendar, [it is easier to create a service account](https://console.cloud.google.com/iam-admin/serviceaccounts){:target="_blank"}. This will give you an email address associated with the service account. If you make sure that account is added to the recurring dog walking meetings then you can programmatically execute this without having to constantly refresh. Alternatively, if you have a personal Google Workspace, you can use domain wide delegation. Sadly that is not an option for us free Gmail users. 

Once you have created the service account in a project that you have enabled read only Calendar API access to, download the JSON credentials it gives you and stash it in the config directory. The app should pull it from there directly, but TBD as more dev happens.