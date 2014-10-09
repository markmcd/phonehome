var app = angular.module('phonehome', []);
app.config(function($sceDelegateProvider) {
    $sceDelegateProvider.resourceUrlWhitelist([
        'self',
        'http://gist-it.appspot.com/**'
    ]);
});
app.controller('searchCtrl', function($scope) {
    $scope.query = '';
    $scope.search = function() {
        var results = search($scope.query);
        $scope.results = results;
        window.setTimeout(function(){
            Prism.highlightAll();
        }, 10);
    }
    $scope.gitHubUrl = function(result) {
        return 'http://gist-it.appspot.com/github/' + result.User + '/' +
            result.Repo + '/blob/master/' + result.Path + '?slice=' +
            (result.Line - 3) + ':' + (result.Line + 3)
    }
});

function search(q) {
    var r;

    $.ajax({
        url: "search?q=" + encodeURIComponent(q),
        dataType: 'json',
        success: function(results) {
            if (results == null) {
                // TODO: handle errors, etc
                return;
            }

            r = results;
        },
        async: false
    });

    return r;
}
